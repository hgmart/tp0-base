import socket
import logging
from common.utils import *

class Server:

    def __init__(self, port, listen_backlog, protocol_payload):
        self.is_enabled = True
        self._server_socket = None
        self.port = port
        self.listen_backlog = listen_backlog
        self.protocol_payload = protocol_payload

    def __enter__(self):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', self.port))
        self._server_socket.listen(self.listen_backlog)
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        if ( self._server_socket is not None):
            self._server_socket.close()
        return True

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        
        while self.is_enabled:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except Exception as e:
                logging.error(f"action: server_failure | result: fail | detail: error raised while waiting new incoming connection: {e}")

    def close(self, signal_number, frame):
        logging.debug(f"action: sigterm_signal | result: success | detail: signal {signal_number} received")
        self.is_enabled = False
        self._server_socket.close()

    def __recursive_receive(self, client_sock):

        bytes = client_sock.recv(self.protocol_payload)
        logging.info(f'action: msg_received | bytes: {len(bytes)}')

        if (len(bytes) < self.protocol_payload):
            return bytes
        
        more_bytes = self.__recursive_receive(client_sock)

        return bytes + more_bytes
    
    def __process_message(self, bytes):
        chunk = bytes.split(b'\x00')
        message_type = chunk.pop(0).decode('utf-8')
        information = [ data.decode('utf-8') for data in chunk[:len(chunk)-1]]
        logging.debug(information)
        
        if (message_type == 'S'):    
            bet = Bet("0", information[0], information[1], information[2], information[3], information[4])
            store_bets([bet])
            return True
        
        return False

    def __recursive_send(self, client_sock, data):

        chunkSize = len(data) if len(data) < self.protocol_payload else self.protocol_payload
        chunk = data[:chunkSize]

        bytes = client_sock.send(chunk)

        if ( bytes == len(data) ):
            return True
        
        rec_bytes = self.__recursive_send(client_sock, data[bytes:])

        return bytes + rec_bytes == len(data)

    def __handle_received_message(self, client_sock, bytes):
        try:
            succeeded = self.__process_message(bytes)

            if succeeded:            
                status_code = "200"
                logging.info("action: processing_message | result: success | status_code: 200 | msg: bet saved")
                self.__recursive_send(client_sock, status_code.encode('utf-8'))

            else:
                status_code = "400"
                logging.error("action: processing_message | result: fail | status_code: 400 | msg: incorrect format")
                self.__recursive_send(client_sock, status_code.encode('utf-8'))

        except Exception as e:
            status_code = "422"
            logging.error(f"action: processing_message | result: fail | status_code: 422 | error: {e}")
            self.__recursive_send(client_sock, status_code.encode('utf-8'))

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            bytes = self.__recursive_receive(client_sock)
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {bytes}')
            self.__handle_received_message(client_sock, bytes)

        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

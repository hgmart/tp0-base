import socket
import logging


class Server:

    def __init__(self, port, listen_backlog):
        self.is_enabled = True
        self._server_socket = None
        self.port = port
        self.listen_backlog = listen_backlog

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
                logging.error(f"action: server_failure | result: fail | error raised while waiting new incoming connection: {e}")

    def close(self, signal_number, frame):
        logging.debug(f"action: sigterm_signal | result: success | signal {signal_number} received")
        self.is_enabled = False
        self._server_socket.close()

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # TODO: Modify the receive to avoid short-reads
            msg = client_sock.recv(1024).rstrip().decode('utf-8')
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # TODO: Modify the send to avoid short-writes
            client_sock.send("{}\n".format(msg).encode('utf-8'))
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

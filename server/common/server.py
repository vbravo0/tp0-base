import socket
import logging
import signal
from common import communication
from common import utils

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        
        # Graceful exit
        self.is_running = True
        signal.signal(signal.SIGTERM, self.exit_gracefully)
            
    def exit_gracefully(self, *args):
        signal = args[0]
        logging.info(f'action: exit_gracefully | signal: {signal}')
        self.is_running = False

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while self.is_running:
            client_sock = self.__accept_new_connection()
            self.__handle_client_connection(client_sock)

        self._server_socket.close()

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            bet = communication.recv_bet(client_sock)
            communication.send_string(client_sock, 'action: apuesta_recibida | result: success')
            utils.store_bets([bet])
            logging.info(f'action: apuesta_almacenada | result: success | dni {bet.document} | numero: {bet.number}')
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

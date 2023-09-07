import socket
import logging
import signal
from common import communication
from common import utils
from common import bet_serializer
from multiprocessing import Process, Lock, Pool, Manager

STORE_SUCCESS = True 
STORE_ERROR = False
WINNERS_SUCCESS = True 
WINNERS_ERROR = False

def handle_store_bets(args):
    """
    Read message from a specific client socket and closes the socket

    If a problem arises in the communication with the client, the
    client socket will also be closed
    """
    client_sock = args[0]
    lock = args[1]

    while True:
        try:
            chunk = communication.recv_string(client_sock)
            if len(chunk) == 0:
                communication.send_string(client_sock, "ok")
                return STORE_SUCCESS
            bets = bet_serializer.bets_from_chunk(chunk)
            lock.acquire()
            utils.store_bets(bets)
            lock.release()
            logging.info(f'action: apuesta_almacenada | result: success')
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
            return STORE_ERROR
    
def handle_winners(args):
    client_sock = args[0]
    winners = args[1]

    agency = communication.recv_u32(client_sock)
    agency_winners = winners.get(agency, [])
    chunk = bet_serializer.bet_documents_to_chunk(agency_winners)
    communication.send_string(client_sock, chunk)
    client_sock.close()
    return WINNERS_SUCCESS

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.agencies = listen_backlog
        
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
        manager = Manager()
        with Pool(processes=self.agencies) as pool:
            while self.is_running:
                client_socks = []
                lock = manager.Lock()

                for _ in range(self.agencies):
                    client_sock = self.__accept_new_connection()
                    client_socks.append(client_sock)

                ret = pool.map(
                    func=handle_store_bets, 
                    iterable=[(client_sock, lock) for client_sock in client_socks]
                )
                logging.info(f'store_bet ret: {ret}')

                winners = self.check_winners()

                ret = pool.map(
                    func=handle_winners, 
                    iterable=[(client_sock, winners) for client_sock in client_socks]
                )
                logging.info(f'winner ret: {ret}')

                        
        self._server_socket.close()
    
    def check_winners(self):
        bets = utils.load_bets()
        winners = [bet for bet in bets if utils.has_won(bet)]

        res = {}
        for winner in winners:
            if winner.agency not in res.keys():
                res[winner.agency] = []
            res[winner.agency].append(winner)

        logging.info(f"res: {res}")
        return res

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

      


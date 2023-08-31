BUFSIZE = 1024

def send(socket, msg: bytes):
    total_sent = 0
    
    while total_sent < len(msg):
        sent = socket.send(msg[total_sent:]) 
        if send == 0:
            raise RuntimeError("socket send error")
        total_sent += sent 

def recv(socket, size: int) -> bytes:
    data_received = b''

    while len(data_received) < size:
        chunk = socket.recv(min(size - len(data_received), BUFSIZE))
        if len(chunk) == 0:
            raise RuntimeError("socket recv error")
        data_received += chunk
    return data_received
    


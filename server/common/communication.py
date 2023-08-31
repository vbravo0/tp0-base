BUFSIZE = 1024
SIZE_I32 = 4
ENDIAN_ORDER = 'big'
STRING_ENCODING = 'ascii'

def send_bytes(socket, msg: bytes):
    total_sent = 0
    while total_sent < len(msg):
        sent = socket.send(msg[total_sent:]) 
        if sent == 0:
            raise RuntimeError("socket send error")
        total_sent += sent 

def recv_bytes(socket, size: int) -> bytes:
    data_received = b''
    while len(data_received) < size:
        chunk = socket.recv(min(size - len(data_received), BUFSIZE))
        if len(chunk) == 0:
            raise RuntimeError("socket recv error")
        data_received += chunk
    return data_received

def send_i32(socket, n: int):
    data = n.to_bytes(SIZE_I32, byteorder=ENDIAN_ORDER, signed=True)
    send_bytes(socket, data)

def recv_i32(socket) -> int:
    data = recv_bytes(socket, SIZE_I32)
    return int.from_bytes(data, byteorder=ENDIAN_ORDER, signed=True)

def send_string(socket, s: str):
    data = bytes(s, encoding=STRING_ENCODING)
    send_bytes(socket, data)

def recv_string(socket, size: int) -> str:
    data = recv_bytes(socket, size)
    return data.decode(str=STRING_ENCODING)



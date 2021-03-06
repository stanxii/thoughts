import grpc
from concurrent import futures
import time
import logging

from authservice import thoughts_pb2_grpc
from authservice.auth import AuthService
from authservice.session import SessionService
from authservice.db import Database
from authservice.db_client import DbClient


class Server:
    def __init__(self):
        self.config = {}
        self.db_client = None

    def get_db_client(self):
        if self.db_client is None:
            db = Database(self.config['DB_URI'])
            self.db_client = DbClient(db)
        return self.db_client

    def create_server(self):
        db_client = self.get_db_client()
        secret = self.config['SECRET']

        auth_service = AuthService(db_client, secret)
        session_service = SessionService(db_client, secret)

        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        thoughts_pb2_grpc.add_AuthServiceServicer_to_server(auth_service, server)
        thoughts_pb2_grpc.add_SessionServiceServicer_to_server(session_service, server)

        server.add_insecure_port(f'[::]:{self.config["PORT"]}')

        return server

    def serve(self):
        server = self.create_server()
        server.start()
        logging.info(f'Server running on port {self.config["PORT"]}')
        try:
            while True:
                time.sleep(60 * 60 * 24)
        except KeyboardInterrupt:
            server.stop(0)

import os

from authservice.server import Server


def create_app():
    app = Server(5010)
    app.config['DB_URI'] = os.getenv('DB_URI')
    app.config['SECRET'] = os.getenv('SECRET')
    return app

import os

from authservice import db
from authservice.server import Server


def create_app():
    app = Server()
    app.config['DB_URI'] = os.getenv('DB_URI')
    app.config['JWT_SECRET'] = os.getenv('JWT_SECRET')
    app.config['PORT'] = os.getenv('PORT')
    return app

from authservice import thoughts_pb2, thoughts_pb2_grpc
from authservice.utils import validate_email, object_to_session
from authservice.auth import AuthException


class AuthService(thoughts_pb2_grpc.AuthServiceServicer):
    def __init__(self, db_client, secret):
        self.db_client = db_client
        self.secret = secret

    def GetSessions(self, request, context):
        """Get all active sessions of the active user."""

        try:
            payload = self.validate_session(request.token)
        except AuthException as e:
            error = thoughts_pb2.Error(code=e.code, error=e.error,
                message=e.message)
            return thoughts_pb2.Status(error=error)

        sessions = self.db_client.get_user_sessions(payload['sub'])
        sessions = [object_to_session(item) for item in sessions]

        return thoughts_pb2.Sessions(sessions=sessions)

    def DeleteSession(self, request, context):
        """Delete a session with a given token."""

        try:
            payload = self.validate_session(request.token)
        except AuthException as e:
            error = thoughts_pb2.Error(code=e.code, error=e.error,
                message=e.message)
            return thoughts_pb2.Status(error=error)

        session_id = request.session_id

        if self.db_client.get_session(session_id)['user_id'] != payload['sub']:
            error = thoughts_pb2.Error(code=403, error='FORBIDDEN',
                message='This action is forbidden.')
            return thoughts_pb2.Status(error=error)

        self.db_client.delete_session(session_id)

        return thoughts_pb2.Status(message='Deleted session.')

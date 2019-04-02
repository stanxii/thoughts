from flask import (
    Blueprint, request, make_response, jsonify
)

from userservice import db


bp = Blueprint('follow', __name__)

@bp.route('/users/<username>/followers', methods=['GET'])
def get_followers(username):
    """Returns users followed by the user."""
    return make_response(jsonify({'response': 'Get followers'}), 200)


@bp.route('/users/<username>/following', methods=['GET'])
def get_following(username):
    """Returns users following the user."""
    return make_response(jsonify({'response': 'Get following'}), 200)


@bp.route('/users/<username>/following', methods=['POST'])
def follow_user(username):
    """Follows a user with the current user."""
    return make_response(jsonify({'response': 'Follow user'}), 200)


@bp.route('/users/<username>/following', methods=['DELETE'])
def unfollow_user(username):
    """Unfollows a user with the current user."""
    return make_response(jsonify({'response': 'Unfollow user'}), 200)

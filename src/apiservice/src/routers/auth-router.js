import { APIRouter } from './api-router';
import { getToken } from './utils';

export class AuthRouter extends APIRouter {
  constructor(authClient) {
    super();
    this.authClient = authClient;
  }

  connectRouter(router) {
    router.get('/', (req, res) => {
      this.getSessions(req, res);
    });

    router.delete('/:id', (req, res) => {
      this.deleteSession(req, res);
    });

    router.post('/', (req, res) => {
      if (req.body.email !== undefined) {
        this.login(req, res);
      } else {
        this.refreshToken(req, res);
      }
    });
  }

  // Authentication

  login(req, res) {
    let email = req.body.email;
    let password = req.body.password;
    let userAgent = req.headers['User-Agent'];

    this.handleResponse(
      this.authClient.createSession(email, password, userAgent), res
    );
  }

  refreshToken(req, res) {
    let token = getToken(req);

    this.handleResponse(this.authClient.refreshToken(token), res);
  }

  // Sessions

  getSessions(req, res) {
    let token = getToken(req);

    this.handleResponse(this.authClient.getSessions(token), res);
  }

  deleteSession(req, res) {
    let token = getToken(req);
    let sessionId = req.params.id;

    this.handleResponse(this.authClient.deleteSession(sessionId, token), res);
  }
}

import { Server } from './server';

const port = Number(process.env.PORT);
const apiRoot = process.env.API_ROOT;

const server = new Server(port, apiRoot);
server.config['AUTH_URI'] = process.env.AUTH_URI;
server.config['POST_URI'] = process.env.POST_URI;
server.config['USER_URI'] = process.env.USER_URI;
server.config['IMAGE_URI'] = process.env.IMAGE_URI;

if (!module.parent) {
  server.start();
}

import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import './App.css';
import Chat from './components/Chat';
import Conversations from './components/Conversations';
import Login from './components/Login';
import Register from './components/Register';
import MessagesContext from './context/messages-context';
import WebSocketContext from './context/socket-context';

function App() {
  const webSocket = null;//new WebSocket(`ws://${window.location.hostname}:3000/ws`);
  const [socket, setSocket] = React.useState(webSocket);
  const [messages, setMessages] = React.useState(MessagesContext);

  return (
    <div className="App">
      <WebSocketContext.Provider value={{socket, setSocket}}>
        <MessagesContext.Provider value={{messages, setMessages}}>
          <Switch>
            <Route path='/' exact component={Login} />
            <Route path='/login' component={Login} />
            <Route path='/register' component={Register} />
            <Route path='/chat/:user' component={Chat} />
            <Route path='/conversations' component={Conversations} />
          </Switch>
        </MessagesContext.Provider>
      </WebSocketContext.Provider>
    </div>
  );
}

export default withRouter(App);

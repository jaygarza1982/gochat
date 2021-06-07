import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import './App.css';
import Chat from './components/Chat';
import WebSocketContext from './context/socket-context';

function App() {
  const webSocket = new WebSocket(`ws://${window.location.hostname}:3000/ws`);
  const [socket] = React.useState(webSocket);

  return (
    <div className="App">
      <WebSocketContext.Provider value={{socket}}>
        <Switch>
          <Route path='/chat' component={Chat} />
        </Switch>
      </WebSocketContext.Provider>
    </div>
  );
}

export default withRouter(App);

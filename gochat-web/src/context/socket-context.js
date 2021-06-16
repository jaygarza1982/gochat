import React from 'react';

const WebSocketContext = React.createContext({
    socket: null,
    setSocket: () => {},
});

export default WebSocketContext;
import React from 'react';

const MessagesContext = React.createContext({
    messages: [],
    setMessages: () => {},
});

export default MessagesContext;
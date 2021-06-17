import React, { useEffect } from 'react';
import { useState } from 'react';
import { useContext } from 'react';
import WebSocketContext from '../../context/socket-context';

const Chat = () => {
    const { socket, setSocket } = useContext(WebSocketContext);

    const [messages, setMessages] = useState([]);

    const handleKeyDown = event => {
        // {"SenderId":"2","ReceiverId":"2","MessageText":"test"} 
        if (event.key === 'Enter' && event.target.value != '') {
            socket.send(event.target.value);
            event.target.value = '';
        }
    }

    const onMessage = event => {
        setMessages([...messages, event.data]);
    }

    // When our socket is set, setup functions
    useEffect(() => {
        if (!socket) return;

        console.log('Setting up socket message handler');
    
        socket.onmessage = onMessage;
    }, [socket, messages]);

    useEffect(() => {
        fetch('/api/test').then(resp => {
            return resp.text();
        }).then(text => {
            console.log(text);
        })
    }, []);

    return (
        <>
            <div className="messages">
                {
                    messages.map((message, index) => {
                        return (
                            <div
                                className="message-secondary secondary-message-color"
                            >
                                {message}
                            </div>
                        )
                    })
                }
                {/* <div className="w-3/5 mx-4 my-2 p-2 rounded-lg secondary-message">Message from other</div>
                <div className="w-3/5 mx-4 my-2 p-2 rounded-lg primary-message-color float-right" style={{marginBottom: 70}}>Message from me</div> */}
            </div>

            <div className="fixed w-full flex justify-between background-color" style={{bottom: 0}}>
                <input
                    onKeyDown={handleKeyDown}
                    style={{margin: 10}}
                    type="text"
                    placeholder="Message..."
                    className="text-box"
                />
            </div>
        </>
    );
}

export default Chat;
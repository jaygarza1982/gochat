import axios from 'axios';
import React, { useEffect } from 'react';
import { useState } from 'react';
import { useContext } from 'react';
import { useParams } from 'react-router-dom';
import WebSocketContext from '../../context/socket-context';

const Chat = () => {
    const { socket, setSocket } = useContext(WebSocketContext);

    const [receiver, setReceiver] = useState(useParams()?.user);
    const [messages, setMessages] = useState([]);

    const scrollToBottom = () => {
        window.scrollTo(0, document.body.scrollHeight);
    }

    const handleKeyDown = event => {
        if (event.key === 'Enter' && event.target.value != '') {
            const postData = { ReceiverId: receiver || '', MessageText: event.target.value };
            axios.post('/api/send-message', postData).then(resp => {
                console.log('Message sent successfully!');

                setMessages([...messages, JSON.stringify(postData)]);

                event.target.value = '';
            }).catch(err => {
                console.log('Could not send user message', err);
            });
        }
    }

    const onMessage = event => {
        // Get message ID and request message from API
        (async () => {
            try {
                const message = await axios.post('/api/get-message', { id: parseInt(event.data) });

                console.log('message received', message.data);
                setMessages([...messages, JSON.stringify(message.data)]);

                scrollToBottom();
            } catch (err) {
                console.log('Could not get message', err);
            }
        })();
    }

    // TODO: disconnect when umount

    // When our socket is set, setup functions
    useEffect(() => {
        if (!socket) return;

        console.log('Setting up socket message handler');
    
        socket.onmessage = onMessage;
    }, [socket, messages]);

    useEffect(() => {
        //Check if we are using http or https
        const https = window.location.protocol.startsWith('https');

        //Setup our socket on load
        const webSocketURL = `${https ? 'wss' : 'ws'}://${window.location.host}/ws`;
        console.log('Setting up websocket at', webSocketURL)
        setSocket(new WebSocket(webSocketURL));

        //Load messages
        (async () => {
            try {
                const messagesFetch = await axios.post('/api/list-messages', { username: receiver });
                const messagesData = messagesFetch.data;

                setMessages(messagesData.map(m => JSON.stringify(m)));

                scrollToBottom();
            } catch (err) {
                console.log('Unable to load messages', err);
            }
        })();
    }, []);

    return (
        <>
            <div className="chat">
                <div className="messages">
                    {
                        messages.map((message, index) => {
                            return (
                                <div className={`message ${JSON.parse(message)?.SenderId == receiver ? 'yours' : 'mine'}`} key={`${index}-message`}>
                                    {JSON.parse(message)?.MessageText}
                                </div>
                            )
                        })
                    }
                </div>
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
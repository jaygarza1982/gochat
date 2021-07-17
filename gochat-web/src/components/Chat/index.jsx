import axios from 'axios';
import React, { useEffect } from 'react';
import { useState } from 'react';
import { useContext } from 'react';
import { useParams } from 'react-router-dom';
import WebSocketContext from '../../context/socket-context';

const Chat = () => {
    const { socket, setSocket } = useContext(WebSocketContext);

    const [receiver, setReceiver] = useState(useParams()?.user);
    const [yourMessages, setYourMessages] = useState([]);
    const [myMessages, setMyMessages] = useState([]);

    const handleKeyDown = event => {
        //TODO: Append to myMessages
        if (event.key === 'Enter' && event.target.value != '') {
            axios.post('/api/send-message', {ReceiverId: receiver || '', MessageText: event.target.value}).then(resp => {
                console.log('Message sent successfully!');

                event.target.value = '';
            }).catch(err => {
                console.log('Could not send user message', err);
            });
        }
    }

    const onMessage = event => {
        const { SenderId } = JSON.parse(event.data);
        
        //TODO: Only append to their messages
        SenderId == receiver ? setMyMessages([...myMessages, event.data]) : setYourMessages([...yourMessages, event.data]);
    }

    // When our socket is set, setup functions
    useEffect(() => {
        if (!socket) return;

        console.log('Setting up socket message handler');
    
        socket.onmessage = onMessage;
    }, [socket, yourMessages, myMessages]);

    useEffect(() => {
        //Setup our socket on load
        setSocket(new WebSocket(`ws://${window.location.hostname}:3000/ws`));
    }, []);

    return (
        <>
            <div className="chat">
                <div className="messages yours">
                    {
                        yourMessages.map((message, index) => {
                            return (
                                <div className="message" key={`${index}-your`}>
                                    {JSON.parse(message)?.MessageText}
                                </div>
                            )
                        })
                    }
                </div>
                <div className="messages mine">
                    {
                        myMessages.map((message, index) => {
                            return (
                                <div className="message" key={`${index}-mine`}>
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
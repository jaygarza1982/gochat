import React, { useEffect } from 'react';
import { useState } from 'react';
import { useContext } from 'react';
import WebSocketContext from '../../context/socket-context';

const Chat = () => {
    const { socket, setSocket } = useContext(WebSocketContext);

    const [messages, setMessages] = useState([]);

    socket.onopen = () => {
        console.log("Successfully Connected");
    };

    socket.onmessage = event => {
        console.log(event.data);

        setMessages([...messages, event.data]);
    };

    const handleKeyDown = event => {
        if (event.key === 'Enter') {
            socket.send(event.target.value);
            event.target.value = '';
        }
    }

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
                    messages.map(message => {
                        return (
                            <div
                            className="w-3/5 mx-4 my-2 p-2
                            rounded-lg secondary-message"
                            >
                                {message}
                            </div>
                        )
                    })
                }
                {/* <div className="w-3/5 mx-4 my-2 p-2 rounded-lg secondary-message">Message from other</div>
                <div className="w-3/5 mx-4 my-2 p-2 rounded-lg primary-message float-right" style={{marginBottom: 70}}>Message from me</div> */}
            </div>

            <div className="fixed w-full flex justify-between background-color" style={{bottom: 0}}>
                <input
                    onKeyDown={handleKeyDown}
                    style={{margin: 10}}
                    type="text"
                    placeholder="Message..."
                    className="px-3 py-3 placeholder-blueGray-300
                    text-blueGray-600 relative
                    bg-white bg-white
                    rounded text-sm border-0
                    shadow outline-none focus:outline-none
                    focus:ring w-full"
                />
            </div>
        </>
    );
}

export default Chat;
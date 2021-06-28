import axios from 'axios';
import React from 'react';
import { useContext } from 'react';
import { withRouter } from 'react-router-dom';
import { useState } from 'react';
import WebSocketContext from '../../context/socket-context';

const Login = props => {

    const { socket, setSocket } = useContext(WebSocketContext);

    const [ form, setForm ] = useState();

    const onTextChange = e => {
        setForm({...form, [e.target.name]: e.target.value});
    }

    const tryLogin = () => {
        axios.post('/api/login', { username: form.goChatUsername, password: form.goChatPassword }).then(resp => {
            console.log(resp);
            if (resp.status === 200) {
                //Set our websocket now that we are authenticated
                setSocket(new WebSocket(`ws://${window.location.hostname}:3000/ws`));

                props.history.push('/chat');
            }
        }).catch(err => {
            console.log('Could not send login', err);
        });
    }

    return (
        <div className="m-3" style={{ color: 'white' }}>
            <div className="w-1/4 m-auto">
                <form className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
                    <div className="mb-4">
                        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="goChatUsername">
                            Username
                        </label>
                        <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                            id="goChatUsername"
                            name="goChatUsername"
                            type="text"
                            placeholder="Username"
                            onChange={onTextChange}
                        />
                    </div>
                    <div className="mb-6">
                        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="goChatPassword">
                            Password
                        </label>
                        <input className="shadow appearance-none border border-red-500 rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
                            id="goChatPassword"
                            name="goChatPassword"
                            type="password"
                            placeholder="*********"
                            onChange={onTextChange}
                        />
                    </div>
                    <div className="flex items-center justify-between">
                        <button
                            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                            type="button"
                            onClick={tryLogin}
                        >
                            Sign In
                        </button>
                    </div>
                </form>
                <p className="text-center text-gray-500 text-xs">
                    GoChat
                </p>
            </div>
        </div>
    );
}

export default withRouter(Login);
import axios from 'axios';
import React from 'react';
import { useContext } from 'react';
import { useEffect } from 'react';
import { withRouter } from 'react-router-dom';
import WebSocketContext from '../../context/socket-context';

const Login = props => {

    const { socket, setSocket } = useContext(WebSocketContext);

    useEffect(() => {
        const { username, password } = { username: 'jake', password: '' };

        axios.post('/api/login', { username, password }).then(resp => {
            console.log(resp);
            if (resp.status === 200) {
                //Set our websocket now that we are authenticated
                setSocket(new WebSocket(`ws://${window.location.hostname}:3000/ws`));
            }
        }).catch(err => {
            console.log('Could not send login', err);
        });
    }, []);

    useEffect(() => {
        props.history.push('/chat');
    }, [socket]);

    return (
        <div style={{color: 'white'}}>
            Login page
            {/* TODO: Add input fields */}
        </div>
    );
}

export default withRouter(Login);
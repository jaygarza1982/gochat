import React, { useEffect } from 'react';

const Chat = () => {
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
                <div className="w-3/5 mx-4 my-2 p-2 rounded-lg secondary-message">Message from other</div>
            </div>
        </>
    );
}

export default Chat;
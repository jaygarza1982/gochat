import React from 'react';
import { withRouter } from 'react-router-dom';

const ConvoEntry = props => {

    const navigateToChat = () => {
        props.history.push(`/chat/${props.name}`);
    }

    return (
        <tr onClick={navigateToChat}>
            <td className="px-6 py-4 whitespace-nowrap">
                <div className="flex items-center">
                    <div className="ml-4">
                        {props.name}
                    </div>
                </div>
            </td>
        </tr>
    );
}

export default withRouter(ConvoEntry);
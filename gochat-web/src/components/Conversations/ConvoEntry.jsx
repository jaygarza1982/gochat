import React from 'react';

const ConvoEntry = ({name}) => {
    return (
        <tr>
            <td className="px-6 py-4 whitespace-nowrap">
                <div className="flex items-center">
                    <div className="ml-4">
                        {name}
                    </div>
                </div>
            </td>
        </tr>
    );
}

export default ConvoEntry;
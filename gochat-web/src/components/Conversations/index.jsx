import React, { useEffect, useState } from 'react';
import ConvoEntry from './ConvoEntry';

const Conversations = () => {

    const [ names, setNames ] = useState(['John Doe', 'Jane Doe']);

    useEffect(() => {
        console.log('names', names);
    });

    return (
        <>
            <div className="flex flex-col">
                <div className="align-middle inline-block min-w-full">
                    <div className="overflow-hidden border-b border-gray-200">
                        <table className="min-w-full divide-y divide-gray-200">
                            <tbody className="bg-white divide-y divide-gray-200">
                                {
                                    names.map(name => {
                                        return <ConvoEntry name={name} />
                                    })
                                }
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </>
    );
}

export default Conversations;
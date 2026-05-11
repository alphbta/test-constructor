import React from 'react';

export default function ShareLinkBox({ link }) {
    const handleCopy = () => {
        navigator.clipboard.writeText(link);
    };
    return (
        <div className="share-link-box">
            <div className="share-link-box-inner">
                <p>Поделиться ссылкой</p>
            </div>
            <div className="share-link-box-inner1">
            <div className="share-link-input-box">
                <input type="text" value={link} readOnly />
                <button onClick={handleCopy}></button>
            </div>
            </div>
        </div>
    );
}

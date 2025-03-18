import React, { useState } from 'react';
import {Tldraw} from "tldraw";
import 'tldraw/tldraw.css'

const Whiteboard = () => {


    return (
        <div style={{ padding: '20px', textAlign: 'center' }}>
            <h2>Интерактивная доска</h2>
            <div style={{height: '800px'}}>
                <Tldraw/>
            </div>
        </div>
    );
};

export default Whiteboard;
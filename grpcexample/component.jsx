import React, { Component, use, useEffect } from 'react';

const comp = (props) => {
    const [state, setState] = React.useState(0);
    useEffect(() => {}, []);

    return <div>{props.text}</div>;
}
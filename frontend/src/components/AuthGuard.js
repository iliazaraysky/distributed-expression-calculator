import React from "react";
import { Redirect } from 'react-router-dom';

const AuthGuard = ( { children }) => {
    const isLoggedIn = checkIfUserIsLoggedIn();

    if  (isLoggedIn) {
        return <>{children}</>
    } else {
        return <Redirect to="/register" />;
    }
};

export default AuthGuard;
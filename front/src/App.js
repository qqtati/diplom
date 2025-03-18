import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import Navbar from './components/Navbar';
import TutorSignUp from './pages/Tutor/SignUp';
import TutorLogin from './pages/Tutor/Login';
import StudentSignUp from './pages/Student/SignUp';
import StudentLogin from './pages/Student/Login';
import MyCalendar from './components/Calendar';
import 'antd/dist/reset.css';
import StudentStatistics from "./pages/StudentStatistics";
import Whiteboard from "./pages/Whiteboard";
import VideoChat from "./pages/VideoChat";
import Profile from "./pages/Profile";
import {AuthProvider} from "./context/AuthContext";

const App = () => {
    return (
        <Router>
            <AuthProvider>
            <Navbar />
            <Switch>
                <Route path="/tutor/signup" component={TutorSignUp} />
                <Route path="/tutor/login" component={TutorLogin} />
                <Route path="/student/signup" component={StudentSignUp} />
                <Route path="/student/login" component={StudentLogin} />
                <Route path="/calendar" component={MyCalendar} />
                <Route path="/tutor/stats" component={StudentStatistics} />
                <Route path="/whiteboard" component={Whiteboard} />
                <Route path="/video" component={VideoChat} />
                <Route path="/profile" component={Profile} />
            </Switch>
            </AuthProvider>
        </Router>

    );
};

export default App;
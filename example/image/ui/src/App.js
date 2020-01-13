import React from "react";
import { Container, Row, Col } from 'react-bootstrap';
import Preview from './components/Preview';

import './App.css'
import 'bootstrap/dist/css/bootstrap.min.css';

const App = props => {
    return (
        <Container className="h-100">
            <Row className="h-100 justify-content-center align-items-center text-center">
                <Col lg={6}>
                    <Preview />
                </Col>
            </Row>
        </Container>
    );
};

export default App;

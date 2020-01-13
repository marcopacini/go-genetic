import {Badge, Card, Image} from "react-bootstrap";
import React, { Component } from "react";
import Toolbar from "./Toolbar";
import {connect} from "react-redux";

class Preview extends Component {
    getURL() {
        if (!this.props.isRunning) {
            return 'http://localhost:3000/placeholder-md.png';
        }

        console.log(this.props.updateCounter);
        return 'http://localhost:3001/best?n=' + this.props.updateCounter;
    }

    render() {
        return (
            <Card>
                <Card.Header>
                    Genetic Art, powered by <i>go-genetic</i>
                </Card.Header>
                <Card.Body>
                    <h5><Badge variant={"secondary"}>n.b. <i>working in progress</i></Badge></h5>
                    <br/><br/>
                    <Image src={this.getURL()} rounded/>
                    <br/><br/>
                    <Toolbar/>
                </Card.Body>
            </Card>
        );
    }
}

export default connect(state => {
    return {
        isRunning: state.isRunning,
        refreshRate: state.refreshRate,
        updateCounter: state.updateCounter
    }
})(Preview);
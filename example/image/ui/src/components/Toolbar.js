import React, { Component } from "react";
import { Button, ButtonGroup, ButtonToolbar, Dropdown, Nav, Navbar, Spinner } from "react-bootstrap";
import { connect } from "react-redux";
import { fetchStats, select, sendStart, update } from "../actions";

class Toolbar extends Component {
    componentDidMount() {
        this.props.dispatch(fetchStats());
    }

    onStart = () => {
        this.props.dispatch(sendStart())
    };

    onUpdate = () => {
        console.log("update");
        this.props.dispatch(update())
    };

    onSelect = (eventKey, event) => {
        this.props.dispatch(select(eventKey));

        if (this.timer != null) {
            clearInterval(this.timer);
        }

        if (eventKey !== '0') {
            const interval = parseInt(eventKey, 10);
            setInterval(this.onUpdate, interval * 1000);
        }
    };

    render() {
        return (
            <Navbar bg="light" expand="md">
                <Navbar.Toggle aria-controls="basic-navbar-nav"/>
                <Nav className="mr-auto">
                    <Run onStart={this.onStart} isRunning={this.props.isRunning}/>
                </Nav>
                <Nav>
                    <Update
                        refreshRate={this.props.refreshRate}
                        onUpdate={this.onUpdate}
                        onSelect={this.onSelect}
                    />
                </Nav>
            </Navbar>
        );
    }
}

const Run = props => {
    return (
        <ButtonToolbar>
            <Button
                variant="outline-primary"
                disabled={props.isRunning}
                onClick={props.onStart}
            >
                {props.isRunning
                    ? <Spinner as="span" animation="grow" size="sm" role="status" aria-hidden="true"/>
                    : "Start"}
            </Button>
        </ButtonToolbar>
    );
};

const REFRESH_RATES = ['1', '5', '30'];

class Update extends Component {
    render() {
        return (
            <Dropdown as={ButtonGroup}>
                <Button
                    variant="outline-secondary"
                    onClick={this.props.onUpdate}
                >
                    Update Now
                </Button>
                <Dropdown.Toggle split variant="outline-secondary" id="dropdown-split-basic" />
                <Dropdown.Menu>
                    <Dropdown.Header>Auto Refresh</Dropdown.Header>
                    {REFRESH_RATES.map(rate => (
                        <Dropdown.Item
                            key={rate}
                            eventKey={rate}
                            active={this.props.refreshRate === rate}
                            onSelect={this.props.onSelect}
                        >
                            {rate} sec
                        </Dropdown.Item>
                    ))}
                    <Dropdown.Divider />
                    <Dropdown.Item
                        eventKey={'0'}
                        active={this.props.refreshRate === '0'}
                        onSelect={this.props.onSelect}
                    >
                        Off
                    </Dropdown.Item>
                </Dropdown.Menu>
            </Dropdown>
        );
    }
}

export default connect(state => {
    return {
        isRunning: state.isRunning,
        refreshRate: state.refreshRate,
        updateCounter: state.updateCounter
    }
})(Toolbar);
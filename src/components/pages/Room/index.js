import { useEffect, useReducer, useContext } from "react";
import { useParams } from "react-router-dom";
import ReconnectingWebSocket from "reconnecting-websocket";

import Grid from "@mui/material/Grid";

import AppContext from "../../contexts";
import reducer from "../../reducers";
import { getRoom } from "../../api";

import Lobby from "./Lobby";
import GameRoom from "./GameRoom";

const Room = () => {
    let { id } = useParams();
    const [state, dispatch] = useReducer(reducer, {
        room: {
            id,
            started: false,
            topic: "",
            users: [],
        },
        you: null,
        joined: false,
    });

    useEffect(() => {
        getRoom({ id }).then((res) => {
            console.log({ res });
            dispatch({
                type: "UPDATE",
                room: res.data,
            });
        });

        const rws = new ReconnectingWebSocket(
            "ws://localhost:5000/api/v1/ws",
            null,
            { reconnectInterval: 10, reconnectDecay: 2 }
        );
        window.addEventListener("beforeunload", () => {
            console.log("beforeunload!");
            rws.close();
        });
        rws.addEventListener("open", () => {
            console.log("Connected");
            rws.send(id);
        });
        rws.addEventListener("message", (msg) => {
            console.log("message recieved: " + msg.data);
            getRoom({ id }).then((res) => {
                console.log({ res });
                dispatch({
                    type: "UPDATE",
                    room: res.data,
                });
            });
        });
        rws.addEventListener("close", () => {
            console.log("connection closed");
        });
        rws.addEventListener("error", (res) => {
            console.log("error");
            console.log(res);
        });

        return () => {
            console.log("component did unmount");
            rws.close();
        };
    }, []);

    return (
        <div style={style}>
            <AppContext.Provider value={{ state, dispatch }}>
                {state.room?.started ? <GameRoom /> : <Lobby />}
            </AppContext.Provider>
        </div>
    );
};

export default Room;

const style = {
    "max-width": "1000px",
    margin: "0 auto",
};

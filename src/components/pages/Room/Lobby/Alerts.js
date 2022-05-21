import React, { useContext } from "react";

import Alert from "@mui/material/Alert";
import AlertTitle from "@mui/material/AlertTitle";
import Stack from "@mui/material/Stack";

import AppContext from "../../../contexts";

const Alerts = () => {
    const { state } = useContext(AppContext);

    let alerts = [];
    if (state.room.users.length === 0) {
        return (
            <Alert key="createUser" severity="info">
                <AlertTitle>Info</AlertTitle>
                プレイヤー名を入力してゲームに参加してください。
            </Alert>
        );
    }
    if (!state.joined) {
        return (
            <Alert key="joinGame" severity="info">
                <AlertTitle>Info</AlertTitle>
                プレイヤー名を入力してゲームに参加するか、自分のプレイヤー名を選択してください。
            </Alert>
        );
    }
    if ([1, 2].includes(state.room.users.length)) {
        alerts.push(
            <Alert key="inviteFriends" severity="info">
                <AlertTitle>Info</AlertTitle>
                URLを共有して友達を招待しましょう。
            </Alert>
        );
    }

    return <Stack>{alerts}</Stack>;
};

export default Alerts;

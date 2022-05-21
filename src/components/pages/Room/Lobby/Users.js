import React, { useContext, useEffect, useState } from "react";

import List from "@mui/material/List";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemText from "@mui/material/ListItemText";
import ListItem from "@mui/material/ListItem";
import IconButton from "@mui/material/IconButton";
import DeleteIcon from "@mui/icons-material/Delete";
import Grid from "@mui/material/Grid";

import AppContext from "../../../contexts";

import { deleteUser } from "../../../api";

const Users = () => {
    const { state } = useContext(AppContext);

    if (typeof state.room === "undefined") {
        return;
    }

    if (state.room.users.length === 0) {
        return;
    }

    return (
        <List>
            {state.room.users.map((user) => {
                return <User key={user.name} user={user} />;
            })}
        </List>
    );
};

export default Users;

const User = ({ user }) => {
    const { state, dispatch } = useContext(AppContext);

    const handleClick = () => {
        dispatch({ type: "ASSIGN_USER", you: user.name });
    };

    const handleDeleteClick = () => {
        deleteUser(state.room.id, user.name);
    };

    return (
        <Grid maxWidth={"sm"}>
            <ListItem>
                <ListItemButton
                    selected={user.isYou}
                    disabled={user.isYou}
                    onClick={handleClick}
                >
                    <PlayerName user={user} />
                </ListItemButton>
                <IconButton
                    onClick={handleDeleteClick}
                    edge="end"
                    aria-label="delete"
                >
                    <DeleteIcon />
                </IconButton>
            </ListItem>
        </Grid>
    );
};

const PlayerName = ({ user }) => {
    if (user.isYou) {
        return <ListItemText primary={user.name} secondary="you" />;
    }
    return <ListItemText primary={user.name} />;
};

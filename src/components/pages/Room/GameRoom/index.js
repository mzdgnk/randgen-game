import React, { useContext } from "react";
import { Controller, useForm } from "react-hook-form";

import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";

import Paper from "@mui/material/Paper";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import CardActions from "@mui/material/CardActions";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Grid from "@mui/material/Grid";
import Box from "@mui/material/Box";
import TextField from "@material-ui/core/TextField";

import AppContext from "../../../contexts";
import { openCard, endGame, startGame } from "../../../api";

const GameRoom = () => {
    const { state, dispatch } = useContext(AppContext);

    return (
        <>
            <Grid container maxWidth="md" justify="center" alignItems="center">
                <Grid item xs={12}>
                    <h1>お題: {state.room.topic}</h1>
                </Grid>
                <Grid item xs={12}>
                    <Users />
                </Grid>
                {/* <PaperUsers /> */}
                <Grid item xs={12}>
                    <StartGameForm />
                </Grid>
            </Grid>
        </>
    );
};

export default GameRoom;

const Users = () => {
    const { state, dispatch } = useContext(AppContext);

    let list = [];
    for (let user of state.room.users) {
        if (user.name === state.you) {
            list.push(<You key={user.name} user={user} />);
        } else {
            list.push(<User key={user.name} user={user} />);
        }
    }

    return (
        <TableContainer component={Paper}>
            <Table aria-label="simple table">
                <TableHead>
                    <TableRow>
                        <TableCell>プレイヤー</TableCell>
                        <TableCell>数字</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {state.room.users.map((user) => {
                        if (user.isYou) {
                            return <You key={user.name} user={user} />;
                        }
                        return <User key={user.name} user={user} />;
                    })}
                </TableBody>
            </Table>
        </TableContainer>
    );
};

const You = ({ user }) => {
    const { state } = useContext(AppContext);

    const handleClick = () => {
        console.log("open card");
        console.log({ state });
        openCard({ roomID: state.room.id, name: user.name });
    };
    return (
        <TableRow key={user.name}>
            <TableCell>{user.name}</TableCell>
            <TableCell>
                {user.num}
                <Button onClick={handleClick}>公開</Button>
            </TableCell>
        </TableRow>
    );
};

const User = ({ user }) => {
    return (
        <TableRow key={user.name}>
            <TableCell>{user.name}</TableCell>
            <TableCell>{user.open ? user.num : "-"}</TableCell>
        </TableRow>
    );
};

const PaperUsers = () => {
    const { state } = useContext(AppContext);

    return (
        <Grid container spacing={2}>
            {state.room.users.map((user) => {
                return <PaperUser key={user.name} user={user} />;
            })}
        </Grid>
    );
};

const PaperUser = ({ user }) => {
    return (
        <Grid item xs={4}>
            <Card sx={{ minWidth: 275 }}>
                <CardContent>
                    <Typography>{user.name}</Typography>
                    {user.isYou ? (
                        <PaperYou user={user} />
                    ) : (
                        <PaperPlayer user={user} />
                    )}
                </CardContent>
            </Card>
        </Grid>
    );
};

const PaperYou = ({ user }) => {
    return (
        <>
            <Typography>{user.num}</Typography>
            <Button size="small">公開</Button>
        </>
    );
};

const PaperPlayer = ({ user }) => {
    if (!user.open) {
        return;
    }
    return <Typography>{user.num}</Typography>;
};

const EndButton = () => {
    const { state, dispatch } = useContext(AppContext);
    const handleClick = () => {
        endGame(state.room.id);
    };

    return (
        <Button color="error" variant="contained" onClick={handleClick}>
            ロビー
        </Button>
    );
};

const StartGameForm = () => {
    const { state, dispatch } = useContext(AppContext);

    const {
        control,
        handleSubmit,
        formState: { errors },
        reset,
    } = useForm({
        defaultValues: {
            topic: "",
        },
    });

    const onSubmit = (event) => {
        startGame({ roomID: state.room.id, topic: event.topic }).then((res) => {
            console.log("started game");
            reset();
        });
    };

    if (!state.joined) {
        return;
    }

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <Grid container spacing={2} alignItems="center" maxWidth="lg">
                <Grid item xs={8}>
                    <Controller
                        control={control}
                        name="topic"
                        rules={{ required: "お題を入力してください。" }}
                        render={({ field }) => (
                            <TextField
                                {...field}
                                label="お題"
                                fullWidth
                                margin="normal"
                                placeholder="お題"
                                error={!!errors.topic}
                                helperText={
                                    errors.topic ? errors.topic.message : ""
                                }
                            />
                        )}
                    />
                </Grid>
                <Grid item xs={2}>
                    <Button
                        variant="contained"
                        color="primary"
                        type="submit"
                        disabled={errors.topic}
                    >
                        次のゲームへ
                    </Button>
                </Grid>
                <Grid item xs={2}>
                    <EndButton />
                </Grid>
            </Grid>
        </form>
    );
};

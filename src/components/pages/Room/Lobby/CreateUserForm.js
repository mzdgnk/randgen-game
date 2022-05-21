import React, { useContext, useState } from "react";
import { useForm, Controller } from "react-hook-form";

import { Button } from "@material-ui/core";
import TextField from "@material-ui/core/TextField";
import Grid from "@mui/material/Grid";

import AppContext from "../../../contexts";

import { createUser } from "../../../api";

const CreateUserForm = () => {
    const { state, dispatch } = useContext(AppContext);
    const [btnDisabled, setBtnDisabled] = useState(false);
    const {
        control,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm({
        defaultValues: {
            name: "",
        },
    });

    const onSubmit = (event) => {
        console.log("create user submitted");
        setBtnDisabled(true);
        if (event.name === "") {
        }
        createUser({ roomID: state.room.id, name: event.name })
            .then((res) => {
                console.log({ createdUsers: res.data });
                dispatch({ type: "ASSIGN_USER", you: res.data.name });
                reset();
            })
            .catch((error) => {
                console.log({ error });
            })
            .finally(setBtnDisabled(false));
    };

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <Grid container spacing={2} alignItems="center" maxWidth="sm">
                <Grid item xs={10}>
                    <Controller
                        control={control}
                        name="name"
                        rules={{
                            required: "プレイヤー名を入力してください。",
                        }}
                        render={({ field }) => (
                            <TextField
                                {...field}
                                label="プレイヤー名"
                                variant="filled"
                                fullWidth
                                margin="normal"
                                error={!!errors.name}
                                helperText={
                                    errors.name ? errors.name?.message : ""
                                }
                                disabled={btnDisabled}
                            />
                        )}
                    />
                </Grid>
                <Grid item xs={2}>
                    <Button
                        variant="contained"
                        color="primary"
                        type="submit"
                        disabled={btnDisabled}
                    >
                        参加
                    </Button>
                </Grid>
            </Grid>
        </form>
    );
};

export default CreateUserForm;

import { useContext } from "react";
import { useForm, Controller } from "react-hook-form";

import { Button } from "@material-ui/core";
import TextField from "@material-ui/core/TextField";
import Grid from "@mui/material/Grid";

import AppContext from "../../../contexts";

import { startGame } from "../../../api";

const StartGameForm = () => {
    const { state, dispatch } = useContext(AppContext);

    const {
        control,
        handleSubmit,
        formState: { errors },
    } = useForm({
        defaultValues: {
            topic: "",
        },
    });

    const onSubmit = (event) => {
        startGame({ roomID: state.room.id, topic: event.topic }).then((res) => {
            console.log("started game");
        });
    };

    if (!state.joined) {
        return;
    }

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <Grid container spacing={2} alignItems="center">
                <Grid item xs={7}>
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
                <Grid item xs={3}>
                    <Button
                        variant="contained"
                        color="primary"
                        type="submit"
                        disabled={errors.topic}
                    >
                        ゲーム開始
                    </Button>
                </Grid>
            </Grid>
        </form>
    );
};

export default StartGameForm;

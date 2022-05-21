import StartGameForm from "./StartGameForm";
import Users from "./Users";
import CreateUserForm from "./CreateUserForm";
import Alerts from "./Alerts";
import Grid from "@mui/material/Grid";

const Lobby = () => {
    return (
        <Grid maxWidth={"md"}>
            <h1>ロビー</h1>
            <Alerts />
            <Users />
            <CreateUserForm />
            <StartGameForm />
        </Grid>
    );
};

export default Lobby;

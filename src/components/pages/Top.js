import { useNavigate } from "react-router-dom";
import { Button } from "@material-ui/core";
import { createRoom as apiCreateRoom } from "../api";
import Grid from "@mui/material/Grid";

const Top = () => {
    const navigate = useNavigate();
    const createRoom = () => {
        apiCreateRoom().then((res) => {
            console.log("data", res);
            let { id } = res.data;
            navigate("/rooms/" + id, { push: true });
        });
    };

    return (
        <div>
            <Grid container>
                <Grid item xs={12}>
                    <h2>
                        プレイヤーに割り当てられたランダムに割り当てられた数字の順番を推測するゲーム
                    </h2>
                </Grid>
                <Grid item xs={12}>
                    <Button
                        variant="contained"
                        color="primary"
                        onClick={createRoom}
                    >
                        ロビー作成
                    </Button>
                </Grid>
            </Grid>
        </div>
    );
};
export default Top;

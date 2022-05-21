import axios from "axios";

const baseURL = process.env.REACT_APP_API_ENDPOINT;

const api = axios.create({
    baseURL: baseURL,
});

export const createRoom = async () => {
    return await api.post("/rooms", {});
};

export const getRoom = async ({ id }) => {
    return await api.get("/rooms/" + id);
};

export const createUser = async ({ roomID, name }) => {
    return await api.post("/rooms/" + roomID + "/users", { name });
};

export const deleteUser = async (roomID, username) => {
    return await api.delete("/rooms/" + roomID + "/users/" + username);
};

export const startGame = async ({ roomID, topic }) => {
    return await api.post("/rooms/" + roomID + "/start", { topic });
};

export const endGame = async (roomID) => {
    return await api.post("/rooms/" + roomID + "/end");
};

export const openCard = async ({ roomID, name }) => {
    return await api.post("/rooms/" + roomID + "/users/" + name + "/open");
};

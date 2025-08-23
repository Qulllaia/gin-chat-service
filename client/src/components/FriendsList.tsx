import { useEffect, useState } from "react";
import FriendCard from "./FriendCard";
import axios from "axios";
import { Friend } from "../types";

export function FriendsList() {

    const [friends, setFriends] = useState<Friend[]>([])

    const fetchFriends = async () => {
        await axios.get('http://localhost:5000/api/user/friends').then((res) => {
            const friendList = res.data.result.map((user: any) => {
                return {
                    id: user.ID,
                    name: user.Name,
                } as Friend;
            })
            setFriends(friendList);
        })
    }

    useEffect(()=> {
        fetchFriends()
    },[])

    return (
        <div className="list">
            {friends.map((friend) => (
                <FriendCard key={friend.id} friend ={friend}/>
            ))}
        </div>
    )
}

export default FriendsList;
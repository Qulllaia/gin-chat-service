import { Friend } from "../types";
import '../styles/FriendCard.css'
interface FriendCardProps {
    friend: Friend
}

export function FriendCard({friend}: FriendCardProps) {
    return (
        <div key={friend.id} className="friend-card">
            <p className="friend-p">{friend.id}</p>
            <p className="friend-p">{friend.name}</p>
        </div>
    )
}

export default FriendCard;
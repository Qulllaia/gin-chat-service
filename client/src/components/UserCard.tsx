import { User } from "../types";
import '../styles/UserCard.css'

interface UserCardProps {
    user: User
    createChatHandler: (user_id: number) => Promise<void>
    handleCheckboxChange: (value: number) => void
    checkboxAvaible: boolean
    isSelected?: boolean
}

export function UserCard({
    user,
    createChatHandler,
    handleCheckboxChange,
    checkboxAvaible,
    isSelected = false,
}: UserCardProps) {
    const checkboxId = `group-user-${user.id}`;

    const handleRowClick = () => {
        if (checkboxAvaible) {
            handleCheckboxChange(user.id);
        } else {
            void createChatHandler(user.id);
        }
    };

    return (
        <div
            className={`user-card${checkboxAvaible && isSelected ? ' user-card--selected' : ''}`}
            onClick={handleRowClick}
            role={checkboxAvaible ? 'checkbox' : undefined}
            aria-checked={checkboxAvaible ? isSelected : undefined}
        >
            {checkboxAvaible && (
                <div
                    className="checkbox"
                    onClick={(e) => e.stopPropagation()}
                >
                    <input
                        className="user-card-checkbox"
                        type="checkbox"
                        id={checkboxId}
                        checked={isSelected}
                        onChange={() => handleCheckboxChange(user.id)}
                    />
                </div>
            )}
            <p className="user-p">{user.name}</p>
        </div>
    )
}

export default UserCard;

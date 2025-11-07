import axios from 'axios';
import { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

export function VerifyPage() {
    const { token } = useParams()
    const navigate = useNavigate()

    const sendVerify = async() => {
        axios.interceptors.request.use(config => {
            config.withCredentials = true;
            return config;
        });
        const result = await axios.get(`http://localhost:5000/api/auth/verify/${token}`)
        .catch((error) => {
            console.log(error);
        })

        if(result?.status === 200) {
            navigate('/chat', { replace: true });
        } 
    }

    useEffect(()=> {
        sendVerify()
    }, [])

    return (
        <div>Verifying</div>
    )
}
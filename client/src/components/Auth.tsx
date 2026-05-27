import axios from "axios"
import { ParentForm } from "./ParentForm"
import { useState } from "react"
import { useNavigate } from "react-router-dom";
import "../styles/Auth.css"


export function Auth() {

    const [name, setName] = useState('');
    const [password, setPassword] = useState('');
    const [email, setEmail] = useState('');
    const [authError, setAuthError] = useState('');
    const [pendingVerification, setPendingVerification] = useState(false);

    const navigate = useNavigate();

    const loginHandler = async () => {
        axios.interceptors.request.use(config => {
            config.withCredentials = true;
            return config;
        });
        const result = await axios.post(`http://localhost:5000/api/auth/login`, { 
            email, 
            name,
            password
        })
        .catch((error) => {
            if(error.status === 500) {
                setAuthError('Fatal Error')
            } else {
                setAuthError(error.response.data.message)
            }
        })

        if(result?.status === 200) {
            navigate('/chat', { replace: true });
        }
    }

    const registrationHandler = async () => {
        axios.interceptors.request.use(config => {
            config.withCredentials = true;
            return config;
        });
        const result = await axios.post(`http://localhost:5000/api/auth/verify`, { 
            email, 
            name,
            password
        })
        .catch((error) => {
            if(error.status === 500) {
                setAuthError('Fatal Error')
            } else {
                setAuthError(error.response.data.message)
            } 
        })

        if (result?.status === 200) {
            setAuthError('');
            setPendingVerification(true);
        }
    }

    return (
        <ParentForm
            isDialog={false}
            isOpen={true}
            setIsOpen={()=>{}}
            backdropClassName="auth-modal-backdrop"
            contentClassName="auth-modal-content"
        >
            <div className="auth-shell">
                {pendingVerification ? (
                    <div className="auth-verify">
                        <div className="auth-verify-icon" aria-hidden>
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                            </svg>
                        </div>
                        <h1 className="auth-verify-title">Проверьте почту</h1>
                        <p className="auth-verify-text">
                            Перейдите по верификационной ссылке, которая придёт на вашу почту.
                            {email && (
                                <>
                                    <br />
                                    <span className="auth-verify-email">{email}</span>
                                </>
                            )}
                        </p>
                        <button
                            className="auth-btn auth-btn-secondary"
                            type="button"
                            onClick={() => setPendingVerification(false)}
                        >
                            Назад ко входу
                        </button>
                    </div>
                ) : (
                    <>
                        <header className="auth-header">
                            <div className="auth-logo" aria-hidden>
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                                </svg>
                            </div>
                            <h1 className="auth-title">Добро пожаловать</h1>
                            <p className="auth-subtitle">Войдите в аккаунт или создайте новый</p>
                        </header>

                        <div className="auth-body">
                            <form onSubmit={(e) => e.preventDefault()}>
                                <div className="auth-field">
                                    <label className="auth-label" htmlFor="auth-email">Email</label>
                                    <div className="auth-input-wrap">
                                        <svg className="auth-input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden>
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                                        </svg>
                                        <input
                                            className="auth-input"
                                            id="auth-email"
                                            type="email"
                                            autoComplete="email"
                                            placeholder="you@example.com"
                                            onChange={(event) => setEmail(event.target.value)}
                                        />
                                    </div>
                                </div>

                                <div className="auth-field">
                                    <label className="auth-label" htmlFor="auth-login">Логин</label>
                                    <div className="auth-input-wrap">
                                        <svg className="auth-input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden>
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                                        </svg>
                                        <input
                                            className="auth-input"
                                            id="auth-login"
                                            type="text"
                                            autoComplete="username"
                                            placeholder="username"
                                            onChange={(event) => setName(event.target.value)}
                                        />
                                    </div>
                                </div>

                                <div className="auth-field">
                                    <label className="auth-label" htmlFor="auth-password">Пароль</label>
                                    <div className="auth-input-wrap">
                                        <svg className="auth-input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden>
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                                        </svg>
                                        <input
                                            className="auth-input"
                                            id="auth-password"
                                            type="password"
                                            autoComplete="current-password"
                                            placeholder="••••••••"
                                            onChange={(event) => setPassword(event.target.value)}
                                        />
                                    </div>
                                </div>

                                <div className="auth-actions">
                                    <button className="auth-btn auth-btn-primary" type="button" onClick={loginHandler}>
                                        Sign in
                                    </button>
                                    <div className="auth-divider">или</div>
                                    <button className="auth-btn auth-btn-secondary" type="button" onClick={registrationHandler}>
                                        Sign up
                                    </button>
                                </div>
                            </form>

                            {authError && (
                                <div className="auth-error-banner" role="alert">
                                    {authError}
                                </div>
                            )}
                        </div>
                    </>
                )}
            </div>
        </ParentForm>
    )
}

export default Auth

import React, { useState, useEffect } from 'react'

export default function App() {
    // ========= auth state =========
    const [token, setToken] = useState('')
    const [isAuthed, setIsAuthed] = useState(false)

    // carrega token salvo (se existir) quando a página abre
    useEffect(() => {
        const saved = localStorage.getItem('hunterToken')
        if (saved && saved !== '') {
            setToken(saved)
            setIsAuthed(true)
        }
    }, [])

    // ========= form login/registro =========
    const [email, setEmail] = useState('')
    const [pass, setPass] = useState('')

    // ========= dados do jogador / feedback UI =========
    const [me, setMe] = useState(null)
    const [msg, setMsg] = useState('')

    // ========= gate state =========
    const [rank, setRank] = useState('C')
    const [minutes, setMinutes] = useState(25)
    const [runId, setRunId] = useState('')
    const [quality, setQuality] = useState(1.0)

    async function register() {
        setMsg('Criando conta...')
        try {
            const res = await fetch('/v1/auth/register', {
                method:'POST',
                headers:{'Content-Type':'application/json'},
                body: JSON.stringify({email, password: pass})
            })

            if (res.status === 409) {
                // e-mail já existe
                setMsg('Essa conta já existe. Só fazer login 👇')
                return
            }

            if (!res.ok) {
                const errTxt = await res.text()
                setMsg('Erro ao registrar: ' + errTxt)
                return
            }

            // se criou com sucesso, tenta logar já
            const loginOk = await login(true)
            if (loginOk) {
                setMsg('Conta criada e logado 🔓')
            } else {
                setMsg('Conta criada! Agora faz login 👇')
            }
        } catch (e) {
            setMsg('Erro ao registrar (fetch falhou)')
        }
    }

    async function login(silent = false) {
        if (!silent) setMsg('Entrando...')
        try {
            const res = await fetch('/v1/auth/login', {
                method:'POST',
                headers:{'Content-Type':'application/json'},
                body: JSON.stringify({email, password: pass})
            })

            // às vezes o backend manda text/plain
            const rawText = await res.text()

            // tenta fazer parse do body na marra
            let data = null
            try {
                data = JSON.parse(rawText)
            } catch {
                // se não deu parse, beleza, fica null
            }

            // tenta pegar token com as duas grafias
            const jwtFromResponse =
                data?.accessToken || // caso certo
                data?.accesToken    // caso com typo (sem o 2º 's')

            if (res.ok && jwtFromResponse) {
                // salva token em memória + localStorage
                setToken(jwtFromResponse)
                localStorage.setItem('hunterToken', jwtFromResponse)

                setIsAuthed(true)

                if (!silent) setMsg('Logado com sucesso 🔓')
                return true
            } else {
                if (!silent) setMsg('Falha no login (credenciais?)')
                return false
            }
        } catch (e) {
            if (!silent) setMsg('Falha no login (fetch)')
            return false
        }
    }


    async function fetchMe() {
        setMsg('Carregando stats...')
        try {
            const res = await fetch('/v1/me', {
                headers:{Authorization:'Bearer '+token}
            })
            const data = await res.json()
            if (res.ok) {
                setMe(data)
                setMsg('Stats atualizados ✅')
            } else {
                setMsg('Erro /me: ' + (data?.error || res.status))
            }
        } catch (e) {
            setMsg('Erro carregando /me (fetch)')
        }
    }

    async function openGate() {
        setMsg('Abrindo gate...')
        try {
            const res = await fetch('/v1/gates/open', {
                method:'POST',
                headers:{
                    'Content-Type':'application/json',
                    Authorization:'Bearer '+token
                },
                body: JSON.stringify({
                    rank,
                    minutes: Number(minutes)
                })
            })
            const data = await res.json()
            if (res.ok && data.id) {
                setRunId(data.id)
                setMsg('Gate aberto. Foca e depois fecha 👊')
            } else {
                setMsg('Erro abrindo gate: ' + (data?.error || res.status))
            }
        } catch (e) {
            setMsg('Erro abrindo gate (fetch)')
        }
    }

    async function closeGate(success=true) {
        if (!runId) {
            setMsg('Nenhum gate ativo pra fechar 😵‍💫')
            return
        }

        setMsg('Fechando gate...')
        try {
            const res = await fetch('/v1/gates/'+runId+'/close', {
                method:'POST',
                headers:{
                    'Content-Type':'application/json',
                    Authorization:'Bearer '+token
                },
                body: JSON.stringify({
                    result: success ? 'success' : 'abandon',
                    quality: Number(quality)
                })
            })
            const data = await res.json()
            if (res.ok && data.result) {
                setMsg(`Gate fechado! XP ganho: ${data.xpEarned} | Gold: ${data.goldEarned}`)
                setRunId('')
            } else {
                setMsg('Erro fechando gate: ' + (data?.error || res.status))
            }
        } catch (e) {
            setMsg('Erro fechando gate (fetch)')
        }
    }

    // ========= render =========

    if (!isAuthed) {
        // tela de login / registro
        return (
            <div className="card">
                <h2 style={{display:'flex',alignItems:'center',gap:'6px'}}>
                    <span>Bem-vindo, hunter</span>
                    <span role="img" aria-label="olho">👁‍🗨</span>
                </h2>

                <label>Email</label>
                <input
                    className="field"
                    value={email}
                    onChange={e=>setEmail(e.target.value)}
                    placeholder="you@example.com"
                />

                <label>Senha</label>
                <input
                    className="field"
                    type="password"
                    value={pass}
                    onChange={e=>setPass(e.target.value)}
                    placeholder="********"
                />

                <button className="btn-main" onClick={()=>login(false)}>Entrar</button>
                <button className="btn-main" onClick={register}>Criar conta</button>

                {msg && <p className="small" style={{marginTop:'12px'}}>{msg}</p>}

                <p className="small" style={{marginTop:'16px', opacity:0.7}}>
                    Ao criar conta você renasce Rank E e começa a farmar XP na marra.
                </p>
            </div>
        )
    }

    // dashboard autenticado
    return (
        <div className="card">
            <h2>Dashboard</h2>

            <button className="btn-main" onClick={fetchMe}>Atualizar stats</button>

            {me && (
                <div style={{marginTop:'12px', fontSize:'14px', lineHeight:'1.4em'}}>
                    <div><strong>Nível:</strong> {me.level}</div>
                    <div><strong>XP:</strong> {me.xp}</div>
                    <div><strong>Gold:</strong> {me.gold}</div>
                    <div><strong>Streak:</strong> {me.streak} 🔥</div>
                </div>
            )}

            <hr className="line" />

            <h3>Abrir Gate (Sessão de Foco)</h3>

            <label>Rank</label>
            <select className="field" value={rank} onChange={e=>setRank(e.target.value)}>
                <option>E</option>
                <option>D</option>
                <option>C</option>
                <option>B</option>
                <option>A</option>
                <option>S</option>
            </select>

            <label>Duração (min)</label>
            <input
                className="field"
                type="number"
                value={minutes}
                onChange={e=>setMinutes(e.target.value)}
            />

            <button className="btn-main" onClick={openGate}>Abrir Gate</button>

            <h4 style={{marginTop:'20px'}}>Fechar Gate</h4>

            <label>Qualidade (0.5 a 1.5)</label>
            <input
                className="field"
                type="number"
                step="0.1"
                value={quality}
                onChange={e=>setQuality(e.target.value)}
            />

            <div className="row">
                <button className="btn-main" onClick={()=>closeGate(true)}>Success ✅</button>
                <button className="btn-main" onClick={()=>closeGate(false)}>Abandon ❌</button>
            </div>

            {msg && (
                <p style={{marginTop:'12px', fontSize:'13px'}}>{msg}</p>
            )}

            {runId && (
                <p className="small" style={{marginTop:'12px', opacity:0.7}}>
                    Gate ativo: {runId}
                </p>
            )}

            <p className="small" style={{marginTop:'24px', opacity:0.5, fontSize:'11px'}}>
                Logado como {email || '(usuário carregado)'}
            </p>
        </div>
    )
}

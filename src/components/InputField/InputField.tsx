import { type FC } from 'react'
import { Button } from 'react-bootstrap'
import './InputField.css'

interface Props {
    value: string
    setValue: (value: string) => void
    onSubmit: () => void
    loading?: boolean
    placeholder?: string
    buttonTitle?: string
}

const InputField: FC<Props> = ({ value, setValue, onSubmit, loading, placeholder, buttonTitle = 'Искать' }) => (
    <div className="search-container">
        <div className="search-bar">
            <input
                value={value}
                placeholder={placeholder}
                onChange={(event => setValue(event.target.value))}
                onKeyDown={(e) => { if (e.key === 'Enter') onSubmit() }}
            />
            <button className="mag_glass" onClick={onSubmit} disabled={loading} aria-label="search">
                <img src="http://127.0.0.1:9000/lb1/search.png" alt="Search" />
            </button>
        </div>
    </div>
)

export default InputField
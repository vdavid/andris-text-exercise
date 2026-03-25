'use client'

import { useRef } from 'react'
import { addMessage } from '@/app/actions'

export default function MessageForm() {
    const formRef = useRef<HTMLFormElement>(null)

    async function handleSubmit(formData: FormData) {
        const result = await addMessage(formData)
        if (result.success) {
            formRef.current?.reset()
        }
    }

    return (
        <form ref={formRef} action={handleSubmit} className="flex gap-3">
            <input
                type="text"
                name="text"
                placeholder="Írj egy üzenetet..."
                required
                className="flex-1 rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20"
            />
            <button
                type="submit"
                className="rounded-lg bg-blue-600 px-5 py-2.5 text-sm font-medium text-white transition-colors hover:bg-blue-700 active:bg-blue-800"
            >
                Küldés
            </button>
        </form>
    )
}

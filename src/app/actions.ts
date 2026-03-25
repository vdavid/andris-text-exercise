'use server'

import { revalidatePath } from 'next/cache'
import { supabase } from '@/lib/supabase'

export async function addMessage(formData: FormData): Promise<{ error?: string; success?: boolean }> {
    const text = formData.get('text') as string

    if (!text || text.trim().length === 0) {
        return { error: 'Az üzenet nem lehet üres.' }
    }

    const { error } = await supabase.from('messages').insert({ text: text.trim() })

    if (error) {
        return { error: 'Nem sikerült menteni az üzenetet.' }
    }

    revalidatePath('/')
    return { success: true }
}

export async function deleteMessage(formData: FormData): Promise<void> {
    const id = formData.get('id') as string

    if (!id) {
        return
    }

    await supabase.from('messages').delete().eq('id', id)

    revalidatePath('/')
}

/**
 * Generates a permanently burned-in watermarked preview image using HTML5 Canvas.
 *
 * Security & Performance Benefits:
 * - The watermark is rendered directly into the JPEG pixels (not CSS overlay).
 * - Inspect element in browser only shows the watermarked image; the raw image is never exposed.
 * - Resizes image to max 1600px dimension for fast loading (~200KB-400KB).
 */
export async function createWatermarkedPreview(
  file: File,
  watermarkText = 'FotoIN'
): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    const objectUrl = URL.createObjectURL(file)

    img.onload = () => {
      URL.revokeObjectURL(objectUrl)

      // 1. Calculate resized dimensions (max 1600px)
      const maxDimension = 1600
      let { width, height } = img

      if (width > maxDimension || height > maxDimension) {
        if (width > height) {
          height = Math.round((height * maxDimension) / width)
          width = maxDimension
        } else {
          width = Math.round((width * maxDimension) / height)
          height = maxDimension
        }
      }

      const canvas = document.createElement('canvas')
      canvas.width = width
      canvas.height = height
      const ctx = canvas.getContext('2d')

      if (!ctx) {
        reject(new Error('Canvas 2D context not available'))
        return
      }

      // 2. Draw scaled base image
      ctx.drawImage(img, 0, 0, width, height)

      // 3. Draw tiled diagonal repeating watermark
      ctx.save()
      const angle = (-26 * Math.PI) / 180
      ctx.rotate(angle)

      const fontSize = Math.max(18, Math.round(width / 34))
      ctx.font = `bold ${fontSize}px sans-serif`
      ctx.fillStyle = 'rgba(255, 255, 255, 0.35)'
      ctx.strokeStyle = 'rgba(0, 0, 0, 0.25)'
      ctx.lineWidth = 1.5

      const stepX = fontSize * 11
      const stepY = fontSize * 5
      const bound = Math.max(width, height) * 2.5

      for (let x = -bound; x < bound; x += stepX) {
        for (let y = -bound; y < bound; y += stepY) {
          ctx.strokeText(watermarkText, x, y)
          ctx.fillText(watermarkText, x, y)
        }
      }
      ctx.restore()

      // 4. Draw center watermark emblem
      ctx.save()
      const centerFontSize = Math.max(24, Math.round(width / 22))
      ctx.font = `900 ${centerFontSize}px sans-serif`
      ctx.textAlign = 'center'
      ctx.textBaseline = 'middle'

      const centerText = watermarkText.toUpperCase()
      ctx.shadowColor = 'rgba(0, 0, 0, 0.6)'
      ctx.shadowBlur = 8
      ctx.fillStyle = 'rgba(255, 255, 255, 0.55)'
      ctx.fillText(centerText, width / 2, height / 2)
      ctx.strokeStyle = 'rgba(0, 0, 0, 0.35)'
      ctx.lineWidth = 2
      ctx.strokeText(centerText, width / 2, height / 2)
      ctx.restore()

      // 5. Export to JPEG Blob
      canvas.toBlob(
        (blob) => {
          if (blob) {
            resolve(blob)
          } else {
            reject(new Error('Failed to create watermarked image blob'))
          }
        },
        'image/jpeg',
        0.82
      )
    }

    img.onerror = (err) => {
      URL.revokeObjectURL(objectUrl)
      reject(err)
    }

    img.src = objectUrl
  })
}

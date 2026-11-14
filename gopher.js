/**
 * @fileoverview Defines the GolangGopher class for rendering a scalable, 
 * vector-style Go gopher graphic onto an HTML5 canvas element using the 2D rendering context.
 * Automatically handles responsive resizing and high-density (Retina) display resolutions.
 */

export class GolangGopher {
    /**
     * Initializes the GolangGopher instance, binds the canvas context,
     * performs an initial draw, and sets up a resize observer for responsive updates.
     * 
     * @param {HTMLCanvasElement} id - The target HTML canvas element to render upon.
     * @throws {Error} Throws an error if the canvas element is invalid or missing.
     */
    constructor(id) {
        // Ensure the provided canvas element is valid before proceeding.
        if (!id || !(id instanceof HTMLCanvasElement)) {
            throw new Error("A valid HTMLCanvasElement must be provided to the GolangGopher constructor.")
        }

        this.canvas = id
        this.context = this.canvas.getContext("2d")

        // Validate that the 2D rendering context was successfully acquired.
        if (!this.context) {
            throw new Error("Failed to acquire the 2D rendering context from the canvas element.")
        }

        // Render the initial scene layout.
        this.drawScene()

        // Observe the canvas for size changes to automatically trigger redraws.
        this.resizeObserver = new ResizeObserver(() => this.drawScene())
        this.resizeObserver.observe(this.canvas)
    }

    /**
     * Renders the complete vector-style Go gopher illustration onto the canvas context,
     * calculating proportional dimensions based on the available canvas width and height.
     * 
     * @param {number} width - The logical width of the canvas rendering area.
     * @param {number} height - The logical height of the canvas rendering area.
     */
    drawGopher(width, height) {
        // Abort rendering if the drawing context is unavailable.
        if (!this.context) {
            return
        }

        // Save the current canvas state before applying transformations and styles.
        this.context.save()

        // Calculate proportional dimensions to keep the gopher responsive.
        const minimumDimension = Math.min(width, height)
        const gopherWidth = minimumDimension * 0.7
        const gopherHeight = gopherWidth * 1.2

        const centerPositionX = width / 2
        const topPositionY = (height - gopherHeight) / 2 + gopherHeight * 0.05

        // Apply a drop shadow effect underneath the gopher body.
        this.context.shadowColor = "rgba(0, 0, 0, 0.4)"
        this.context.shadowBlur = 8
        this.context.shadowOffsetY = 4

        // Create a radial gradient to give the gopher body and head a 3D cylindrical appearance.
        const bodyGradient = this.context.createRadialGradient(
            centerPositionX - gopherWidth * 0.15,
            topPositionY + gopherHeight * 0.2,
            gopherWidth * 0.05,
            centerPositionX,
            topPositionY + gopherHeight * 0.5,
            gopherWidth * 0.65
        )

        // Define color stops for the body gradient.
        bodyGradient.addColorStop(0, "#6AD7E5")
        bodyGradient.addColorStop(0.6, "#00ADD8")
        bodyGradient.addColorStop(1, "#0083A0")

        // Fill the main body shape using rounded corners for a soft look.
        this.context.fillStyle = bodyGradient
        this.context.beginPath()
        this.context.roundRect(
            centerPositionX - gopherWidth / 2,
            topPositionY,
            gopherWidth,
            gopherHeight,
            [gopherWidth * 0.45, gopherWidth * 0.45, gopherWidth * 0.25, gopherWidth * 0.25]
        )
        this.context.fill()

        // Disable shadows so facial features render cleanly without unwanted blurring.
        this.context.shadowColor = "transparent"

        // Calculate and draw the left and right ears.
        const earRadius = gopherWidth * 0.12
        const earOffsetY = topPositionY + gopherHeight * 0.12

        this.context.fillStyle = "#00ADD8"

        // Render the left ear.
        this.context.beginPath()
        this.context.arc(centerPositionX - gopherWidth * 0.42, earOffsetY, earRadius, 0, Math.PI * 2)
        this.context.fill()

        // Render the right ear.
        this.context.beginPath()
        this.context.arc(centerPositionX + gopherWidth * 0.42, earOffsetY, earRadius, 0, Math.PI * 2)
        this.context.fill()

        // Calculate positioning metrics for the eyes.
        const eyeRadius = gopherWidth * 0.22
        const leftEyePositionX = centerPositionX - gopherWidth * 0.21
        const rightEyePositionX = centerPositionX + gopherWidth * 0.21
        const eyePositionY = topPositionY + gopherHeight * 0.26

        // Define the horizontal coordinates for both eyes.
        const eyeXPositions = [leftEyePositionX, rightEyePositionX]

        // Iterate through both eyes to render their outer whites, dark pupils, and specular highlights.
        eyeXPositions.forEach((eyePositionX) => {
            // Create a subtle gradient for the eyeball background.
            const eyeGradient = this.context.createRadialGradient(
                eyePositionX - eyeRadius * 0.2,
                eyePositionY - eyeRadius * 0.2,
                eyeRadius * 0.1,
                eyePositionX,
                eyePositionY,
                eyeRadius
            )

            eyeGradient.addColorStop(0, "#FFFFFF")
            eyeGradient.addColorStop(0.85, "#F3F4F6")
            eyeGradient.addColorStop(1, "#D1D5DB")

            // Draw the outer eyeball circle.
            this.context.beginPath()
            this.context.arc(eyePositionX, eyePositionY, eyeRadius, 0, Math.PI * 2)
            this.context.fillStyle = eyeGradient
            this.context.fill()

            const pupilRadius = eyeRadius * 0.48

            // Draw the dark central pupil.
            this.context.beginPath()
            this.context.arc(eyePositionX, eyePositionY, pupilRadius, 0, Math.PI * 2)
            this.context.fillStyle = "#111827"
            this.context.fill()

            // Draw the white reflection highlight inside the pupil for added life.
            this.context.beginPath()
            this.context.arc(
                eyePositionX - pupilRadius * 0.35,
                eyePositionY - pupilRadius * 0.35,
                pupilRadius * 0.25,
                0,
                Math.PI * 2
            )
            this.context.fillStyle = "rgba(255, 255, 255, 0.95)"
            this.context.fill()
        })

        // Calculate and render the lighter-colored snout area.
        const snoutPositionY = eyePositionY + eyeRadius * 0.72

        this.context.beginPath()
        this.context.ellipse(
            centerPositionX,
            snoutPositionY,
            gopherWidth * 0.2,
            gopherWidth * 0.12,
            0,
            0,
            Math.PI * 2
        )
        this.context.fillStyle = "#F5D6BA"
        this.context.fill()

        // Calculate and render the dark nose patch on the snout.
        const nosePositionY = snoutPositionY - gopherWidth * 0.03

        this.context.beginPath()
        this.context.ellipse(
            centerPositionX,
            nosePositionY,
            gopherWidth * 0.08,
            gopherWidth * 0.05,
            0,
            0,
            Math.PI * 2
        )
        this.context.fillStyle = "#1A1A1A"
        this.context.fill()

        // Calculate metrics for rendering the prominent front teeth.
        const teethPositionY = snoutPositionY + gopherWidth * 0.08
        const toothWidth = gopherWidth * 0.05
        const toothHeight = gopherWidth * 0.07

        this.context.fillStyle = "#FFFFFF"

        // Render the left front tooth.
        this.context.beginPath()
        this.context.roundRect(
            centerPositionX - toothWidth - 1,
            teethPositionY,
            toothWidth,
            toothHeight,
            [0, 0, 3, 3]
        )
        this.context.fill()

        // Render the right front tooth.
        this.context.beginPath()
        this.context.roundRect(
            centerPositionX + 1,
            teethPositionY,
            toothWidth,
            toothHeight,
            [0, 0, 3, 3]
        )
        this.context.fill()

        // Restore the canvas state back to before transformations and styles were applied.
        this.context.restore()
    }

    /**
     * Manages high-DPI (Retina) display scaling, clears the drawing buffer,
     * and triggers the main gopher rendering routine.
     */
    drawScene() {
        // Abort if the canvas or rendering context is missing.
        if (!this.canvas || !this.context) {
            return
        }

        // Determine the device pixel ratio to ensure crisp rendering on high-resolution displays.
        const devicePixelRatio = window.devicePixelRatio || 1
        const displayWidth = this.canvas.clientWidth
        const displayHeight = this.canvas.clientHeight

        // Resize the internal drawing buffer if it does not match the scaled display dimensions.
        if (
            this.canvas.width !== displayWidth * devicePixelRatio ||
            this.canvas.height !== displayHeight * devicePixelRatio
        ) {
            this.canvas.width = displayWidth * devicePixelRatio
            this.canvas.height = displayHeight * devicePixelRatio
        }

        // Save state, scale the context for high-DPI displays, and clear the previous frame.
        this.context.save()
        this.context.scale(devicePixelRatio, devicePixelRatio)
        this.context.clearRect(0, 0, displayWidth, displayHeight)

        // Execute the gopher drawing routine using the current display dimensions.
        this.drawGopher(displayWidth, displayHeight)

        // Restore the original context transformation state.
        this.context.restore()
    }
}
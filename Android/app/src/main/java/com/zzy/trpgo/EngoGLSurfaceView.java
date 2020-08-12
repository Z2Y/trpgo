package com.zzy.trpgo;

import android.content.Context;
import android.content.res.Resources;
import android.opengl.GLSurfaceView;
import android.util.Log;
import android.util.AttributeSet;
import android.view.MotionEvent;

import javax.microedition.khronos.egl.EGLConfig;
import javax.microedition.khronos.opengles.GL10;

import bridge.*;

public class EngoGLSurfaceView extends GLSurfaceView {

    private class EngoRenderer implements Renderer {

        private boolean mErrored;

        @Override
        public void onDrawFrame(GL10 gl) {
            if (mErrored) {
                return;
            }
            try {
                Bridge.update();
            } catch (Exception e) {
                Log.e("Go Error", e.toString());
                mErrored = true;
            }
        }

        @Override
        public void onSurfaceCreated(GL10 gl, EGLConfig config) {
        }

        @Override
        public void onSurfaceChanged(GL10 gl, int width, int height) {
        }
    }

    public EngoGLSurfaceView(Context context) {
        super(context);
        initialize();
    }

    public EngoGLSurfaceView(Context context, AttributeSet attrs) {
        super(context, attrs);
        initialize();
    }

    private void initialize() {
        setEGLContextClientVersion(2);
        setEGLConfigChooser(8, 8, 8, 8, 0, 0);
        setRenderer(new EngoRenderer());
    }

    @Override
    public void onLayout(boolean changed, int left, int top, int right, int bottom) {
        super.onLayout(changed, left, top, right, bottom);

        try {
            if (!Bridge.isRunning()) {
                Bridge.start(Resources.getSystem().getDisplayMetrics().widthPixels, Resources.getSystem().getDisplayMetrics().heightPixels);
            }
        } catch (Exception e) {
            Log.e("Go Error", e.toString());
        }
    }

    @Override
    public boolean onTouchEvent(MotionEvent e) {
        for (int i = 0; i < e.getPointerCount(); i++) {
            Bridge.touch((int)e.getX(i), (int)e.getY(i), i, e.getActionMasked());
        }
        return true;
    }

    @Override
    public void onPause() {
        try {
            if (Bridge.isRunning()) {
                // bridge.stop();
            }
        } catch (Exception e) {
            Log.e("Go Error", e.toString());
        }
    }
}